"""Release metadata must agree before the version tag can be pushed."""
from pathlib import Path
import re
import tomllib
import unittest

ROOT = Path(__file__).resolve().parents[2]


class ReleaseMetadataTests(unittest.TestCase):
    def test_root_and_starter_lockfiles_match_package_version(self):
        version = tomllib.loads((ROOT / "Cargo.toml").read_text(encoding="utf-8"))["package"]["version"]
        for relative in ("Cargo.lock", "starters/rust-agent/Cargo.lock"):
            with self.subTest(lockfile=relative):
                packages = tomllib.loads((ROOT / relative).read_text(encoding="utf-8"))["package"]
                package = [p for p in packages if p["name"] == "mote-agent"]
                self.assertEqual(len(package), 1)
                self.assertEqual(package[0]["version"], version)

    def test_current_version_has_nonempty_release_notes(self):
        version = tomllib.loads((ROOT / "Cargo.toml").read_text(encoding="utf-8"))["package"]["version"]
        changelog = (ROOT / "CHANGELOG.md").read_text(encoding="utf-8")
        section = re.search(rf"(?ms)^## {re.escape(version)}\n(.*?)(?=^## |\Z)", changelog)
        self.assertIsNotNone(section, f"Missing changelog section for {version}")
        self.assertTrue(section.group(1).strip(), f"Empty changelog section for {version}")


if __name__ == "__main__":
    unittest.main()
