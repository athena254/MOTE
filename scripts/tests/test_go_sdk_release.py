"""Keep the official Go SDK contract connected to release validation."""
from pathlib import Path
import re
import unittest

ROOT = Path(__file__).resolve().parents[2]


class GoSDKReleaseTests(unittest.TestCase):
    def test_go_sdk_spec_is_in_next_release_not_published_0140(self):
        changelog = (ROOT / "CHANGELOG.md").read_text(encoding="utf-8")
        sections = dict(re.findall(r"(?ms)^## ([^\n]+)\n(.*?)(?=^## |\Z)", changelog))
        self.assertIn("official Go", sections["0.14.1"])
        self.assertIn("GO_SDK_SPEC.md", sections["0.14.1"])
        self.assertNotIn("sdk/go", sections["0.14.0"])
        spec = (ROOT / "docs/GO_SDK_SPEC.md").read_text(encoding="utf-8")
        for contract in ("v0.14.1", "1.22", "synchronously", "ErrTimeout", "1,048,576"):
            self.assertIn(contract, spec)
        readme = (ROOT / "README.md").read_text(encoding="utf-8")
        self.assertIn("docs/GO_SDK_SPEC.md", readme)

    def test_ci_and_extracted_release_exercise_go_with_bridge(self):
        for name in ("ci.yml", "release.yml"):
            workflow = (ROOT / ".github/workflows" / name).read_text(encoding="utf-8")
            with self.subTest(workflow=name):
                self.assertIn("actions/setup-go@", workflow)
                self.assertIn("go-version: '1.22'", workflow)
                self.assertIn("MOTE_BRIDGE", workflow)
                if name == "ci.yml":
                    for command in ("gofmt -l .", "go build ./...", "go vet ./...", "go test ./..."):
                        self.assertIn(command, workflow)
                else:
                    self.assertIn("go -C sdk/go test ./...", workflow)


if __name__ == "__main__":
    unittest.main()
