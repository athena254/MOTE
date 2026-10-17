import json
import os
from pathlib import Path
import queue
import subprocess
import sys
import time

import pytest
from mote_sdk import Client, MoteError, ProtocolError, SetupError, Tool


def fixture(tmp_path, body):
    script = tmp_path / "fixture.py"
    script.write_text("import sys,json,time\nsys.stdin.readline()\n" + body, encoding="utf-8")
    # Allow interpreter startup; dedicated timeout tests keep short explicit deadlines.
    return Client(binary=[sys.executable, str(script)], timeout=15)


@pytest.mark.parametrize("tail, has_extra", [(b"{}\n", True), (None, False)])
def test_delayed_trailing_frame_is_not_ignored(tail, has_extra):
    class PendingOutput(queue.Queue):
        def get(self, *args, **kwargs):
            # Publish on a wait, after any premature empty-queue inspection.
            self.put(tail)
            return super().get(*args, **kwargs)
    assert Client._drain(PendingOutput(), time.monotonic() + 15) is has_extra


def test_no_read_child_cannot_block_startup_forever(tmp_path):
    # Bound this regression from outside so a broken SDK cannot hang pytest.
    sleeping = tmp_path / "sleep.py"
    sleeping.write_text("import time; time.sleep(4)")
    outer = "from mote_sdk import Client,TimeoutError; import sys\ntry:\n Client(binary=[sys.executable,sys.argv[1]],timeout=.2).run({},'x'*900000)\nexcept TimeoutError:\n sys.exit(0)\nsys.exit(1)"
    # The outer guard includes interpreter startup; the SDK deadline stays .2 seconds.
    result = subprocess.run([sys.executable, "-c", outer, str(sleeping)], timeout=15, capture_output=True, env={**os.environ, "PYTHONPATH": str(Path(__file__).parents[1])})
    assert result.returncode == 0, result.stderr.decode()


@pytest.mark.parametrize("timeout", [0, -1, float("inf"), float("nan"), True])
def test_invalid_deadline_rejected(timeout):
    with pytest.raises(SetupError):
        Client(timeout=timeout)


@pytest.mark.parametrize("tail", ["print('{bad}')", "sys.stdout.write('partial')"])
def test_malformed_bytes_after_terminal_fail(tmp_path, tail):
    client = fixture(tmp_path, "print(json.dumps({'type':'run_result','state':'completed','events':[]})); " + tail)
    with pytest.raises(ProtocolError):
        client.run({}, "go")


def test_stderr_flood_does_not_deadlock(tmp_path):
    client = fixture(tmp_path, "sys.stderr.write('x'*200000); print(json.dumps({'type':'run_result','state':'completed','events':[]}))")
    assert client.run({}, "go").success


def test_exit_mismatch_is_not_success(tmp_path):
    client = fixture(tmp_path, "print(json.dumps({'type':'run_result','state':'completed','events':[]})); sys.exit(1)")
    with pytest.raises(MoteError):
        client.run({}, "go")


def test_callback_exception_is_redacted(tmp_path):
    client = fixture(tmp_path, "print(json.dumps({'type':'tool_call','id':1,'name':'x','input':{}}),flush=True)\nr=json.loads(sys.stdin.readline()); assert r=={'type':'tool_result','id':1,'error':'tool failed'}\nprint(json.dumps({'type':'run_result','state':'failed','events':[]})); sys.exit(1)")
    def handler(_):
        raise RuntimeError("PRIVATE CREDENTIAL SHOULD NEVER ENTER PROTOCOL")
    assert not client.run({'capabilities':['x']}, 'go', {'x':Tool('x',handler)}).success
