use serde_json::Value;
use std::cell::RefCell;
use std::io::Read;
use std::time::Duration;

thread_local! {
    // ponytail: one client per thread; timeout changes replace the pool, split by adapter if measured.
    static AGENT: RefCell<Option<(Duration, ureq::Agent)>> = const { RefCell::new(None) };
}

pub(crate) fn post_json(
    endpoint: &str,
    headers: &[(&str, String)],
    body: &Value,
    timeout: Duration,
    max_response: usize,
) -> Result<Value, String> {
    let agent = AGENT.with(|slot| {
        let mut slot = slot.borrow_mut();
        let (_, agent) = slot
            .take()
            .filter(|(cached_timeout, _)| *cached_timeout == timeout)
            .unwrap_or_else(|| {
                (
                    timeout,
                    ureq::AgentBuilder::new()
                        .timeout_connect(timeout)
                        .timeout_read(timeout)
                        .timeout_write(timeout)
                        .build(),
                )
            });
        *slot = Some((timeout, agent.clone()));
        agent
    });
    let mut request = agent
        .post(endpoint)
        .set("Content-Type", "application/json")
        .set("User-Agent", concat!("mote/", env!("CARGO_PKG_VERSION")));
    for (name, value) in headers {
        request = request.set(name, value);
    }
    let response = match request.send_json(body.clone()) {
        Ok(response) => response,
        Err(ureq::Error::Status(status, _)) => return Err(format!("http status {status}")),
        Err(ureq::Error::Transport(error)) => return Err(format!("http transport error: {error}")),
    };
    let mut bytes = Vec::new();
    response
        .into_reader()
        .take((max_response + 1) as u64)
        .read_to_end(&mut bytes)
        .map_err(|error| format!("response read failed: {error}"))?;
    if bytes.len() > max_response {
        return Err("response too large".to_owned());
    }
    serde_json::from_slice(&bytes).map_err(|error| format!("invalid json response: {error}"))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::{BufRead, BufReader, Write};
    use std::net::TcpListener;
    use std::thread;

    #[test]
    fn repeated_requests_reuse_one_connection() {
        let listener = TcpListener::bind("127.0.0.1:0").unwrap();
        let endpoint = format!("http://{}/", listener.local_addr().unwrap());
        let server = thread::spawn(move || {
            let mut requests = 0;
            let mut connections = 0;
            while requests < 4 {
                let (stream, _) = listener.accept().unwrap();
                stream
                    .set_read_timeout(Some(Duration::from_secs(5)))
                    .unwrap();
                let mut stream = BufReader::new(stream);
                connections += 1;
                loop {
                    let mut line = String::new();
                    if stream.read_line(&mut line).unwrap() == 0 {
                        break;
                    }
                    let mut length = 0;
                    loop {
                        line.clear();
                        assert!(stream.read_line(&mut line).unwrap() > 0);
                        if line.as_bytes() == [13, 10] {
                            break;
                        }
                        if let Some((name, value)) = line.split_once(':') {
                            if name.eq_ignore_ascii_case("content-length") {
                                length = value.trim().parse::<usize>().unwrap();
                            }
                        }
                    }
                    let mut body = vec![0; length];
                    stream.read_exact(&mut body).unwrap();
                    stream
                        .get_mut()
                        .write_all(b"HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\n{}")
                        .unwrap();
                    requests += 1;
                    if requests == 4 {
                        return connections;
                    }
                }
            }
            connections
        });
        for _ in 0..4 {
            assert_eq!(
                post_json(
                    &endpoint,
                    &[],
                    &serde_json::json!({}),
                    Duration::from_secs(2),
                    32
                )
                .unwrap(),
                serde_json::json!({})
            );
        }
        assert_eq!(server.join().unwrap(), 1);
    }

    #[test]
    fn cached_client_honors_timeout_changes_and_response_limit() {
        for (timeout, limit, expected) in [
            (Duration::from_secs(5), 32, None),
            (Duration::from_millis(30), 32, Some("http transport error")),
            (Duration::from_secs(5), 1, Some("response too large")),
        ] {
            let listener = TcpListener::bind("127.0.0.1:0").unwrap();
            let endpoint = format!("http://{}/", listener.local_addr().unwrap());
            let server = thread::spawn(move || {
                let (mut stream, _) = listener.accept().unwrap();
                stream
                    .set_read_timeout(Some(Duration::from_secs(5)))
                    .unwrap();
                stream.read_exact(&mut [0]).unwrap();
                thread::sleep(Duration::from_millis(300));
                let _ = stream.write_all(
                    b"HTTP/1.1 200 OK\r\nContent-Length: 2\r\nConnection: close\r\n\r\n{}",
                );
            });
            let result = post_json(&endpoint, &[], &serde_json::json!({}), timeout, limit);
            match expected {
                Some(message) => assert!(result.unwrap_err().contains(message)),
                None => assert_eq!(result.unwrap(), serde_json::json!({})),
            }
            server.join().unwrap();
        }
    }
}
