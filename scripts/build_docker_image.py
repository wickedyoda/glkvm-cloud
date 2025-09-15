#!/usr/bin/env python3
"""Utility for building and publishing the GLKVM Cloud Docker image.

This script replaces the legacy shell-based build helper so that the Docker
image can be produced from Linux, macOS, and Windows alike.  It compiles the
Go binary with the same metadata that the original shell script injected and
then delegates to the local Docker installation to assemble the final image.

Examples
--------
Build a local image tagged as ``latest``::

    python scripts/build_docker_image.py --image ghcr.io/wickedyoda/glkvm-cloud

Build and push images tagged as ``latest`` and ``dev``::

    python scripts/build_docker_image.py \
        --image ghcr.io/wickedyoda/glkvm-cloud \
        --tag latest --tag dev --push

The script assumes that ``go`` and ``docker`` are already installed and
available on ``PATH``.
"""

from __future__ import annotations

import argparse
import datetime as _dt
import os
from pathlib import Path
import shlex
import shutil
import subprocess
import sys
from typing import Iterable, List, Sequence


REPO_ROOT = Path(__file__).resolve().parents[1]


def _format_cmd(cmd: Sequence[object]) -> str:
    """Return a human-friendly representation of *cmd* for logging."""

    return " ".join(shlex.quote(str(part)) for part in cmd)


def _run(cmd: Sequence[object], *, env: dict[str, str] | None = None) -> None:
    """Execute *cmd* and raise :class:`RuntimeError` on failure."""

    print(f"[build] $ {_format_cmd(cmd)}")
    try:
        subprocess.run(cmd, check=True, env=env)
    except subprocess.CalledProcessError as exc:  # pragma: no cover - CLI tool
        raise RuntimeError(
            f"Command failed with exit status {exc.returncode}: {_format_cmd(cmd)}"
        ) from exc


def _ensure_tool(name: str) -> None:
    """Verify that *name* exists on ``PATH``."""

    if shutil.which(name) is None:
        raise RuntimeError(
            f"Required executable '{name}' was not found on PATH. "
            "Please install it before running this script."
        )


def _git(*args: str) -> str:
    """Run a ``git`` command and return its stripped textual output."""

    try:
        return subprocess.check_output(("git",) + args, cwd=REPO_ROOT, text=True).strip()
    except (subprocess.CalledProcessError, FileNotFoundError):
        return "unknown"


def _build_binary(goos: str, goarch: str, output: Path) -> None:
    """Compile the ``rttys`` binary for the requested platform."""

    git_commit = _git("rev-parse", "--short", "HEAD")
    build_time = (
        _dt.datetime.now(_dt.timezone.utc)
        .replace(microsecond=0)
        .strftime("%Y-%m-%dT%H:%M:%S%z")
    )

    env = os.environ.copy()
    env.update({
        "GOOS": goos,
        "GOARCH": goarch,
        "CGO_ENABLED": "0",
    })

    ldflags = (
        f"-s -w "
        f"-X main.GitCommit={git_commit} "
        f"-X main.BuildTime={build_time}"
    )

    # ``go build`` automatically chooses the main module in the repository root.
    cmd: List[object] = [
        "go",
        "build",
        "-trimpath",
        "-ldflags",
        ldflags,
        "-o",
        str(output),
        ".",
    ]

    _run(cmd, env=env)

    # Go on Windows respects the requested output name even if it does not end
    # with ``.exe``; nonetheless we make the binary executable when possible.
    try:
        output.chmod(0o755)
    except OSError:
        pass


def _build_image(
    dockerfile: Path,
    context: Path,
    tags: Iterable[str],
    platform: str | None,
    build_args: Iterable[str],
) -> None:
    """Invoke ``docker build`` using the provided configuration."""

    cmd: List[object] = ["docker", "build", "-f", str(dockerfile)]

    for tag in tags:
        cmd.extend(["-t", tag])

    if platform:
        cmd.extend(["--platform", platform])

    for build_arg in build_args:
        cmd.extend(["--build-arg", build_arg])

    cmd.append(str(context))

    _run(cmd)


def _push(tags: Iterable[str]) -> None:
    for tag in tags:
        _run(["docker", "push", tag])


def parse_args(argv: Sequence[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)

    parser.add_argument(
        "--image",
        default="ghcr.io/wickedyoda/glkvm-cloud",
        help="Base image name (without tag) to use for the Docker build.",
    )
    parser.add_argument(
        "--tag",
        dest="tags",
        action="append",
        default=[],
        help="Tag to apply to the built image (can be passed multiple times).",
    )
    parser.add_argument(
        "--goos",
        default="linux",
        help="GOOS value used when compiling the binary (default: linux).",
    )
    parser.add_argument(
        "--goarch",
        default="amd64",
        help="GOARCH value used when compiling the binary (default: amd64).",
    )
    parser.add_argument(
        "--platform",
        help=(
            "Target platform passed to docker build (for example linux/amd64). "
            "Defaults to the local Docker configuration."
        ),
    )
    parser.add_argument(
        "--dockerfile",
        default="Dockerfile",
        help="Path to the Dockerfile (default: Dockerfile).",
    )
    parser.add_argument(
        "--context",
        default=".",
        help="Docker build context (default: current repository root).",
    )
    parser.add_argument(
        "--build-arg",
        dest="build_args",
        action="append",
        default=[],
        help="Additional build arguments to forward to docker build.",
    )
    parser.add_argument(
        "--push",
        action="store_true",
        help="Push the resulting image tags to the configured registry.",
    )
    parser.add_argument(
        "--keep-binary",
        action="store_true",
        help="Do not delete the compiled binary after the Docker build completes.",
    )

    return parser.parse_args(argv)


def main(argv: Sequence[str] | None = None) -> int:
    args = parse_args(sys.argv[1:] if argv is None else argv)

    # Resolve filesystem locations relative to the repository root.
    dockerfile = (REPO_ROOT / args.dockerfile).resolve()
    context = (REPO_ROOT / args.context).resolve()
    binary_path = REPO_ROOT / "rttys"

    if not dockerfile.exists():
        raise RuntimeError(f"Dockerfile not found: {dockerfile}")

    _ensure_tool("go")
    _ensure_tool("docker")

    if args.tags:
        full_tags = [
            tag if ":" in tag else f"{args.image}:{tag}"
            for tag in args.tags
        ]
    else:
        full_tags = [args.image] if ":" in args.image else [f"{args.image}:latest"]

    print("[build] Compiling Go binary…")
    _build_binary(args.goos, args.goarch, binary_path)

    try:
        print("[build] Building Docker image…")
        _build_image(dockerfile, context, full_tags, args.platform, args.build_args)

        if args.push:
            print("[build] Pushing image tags…")
            _push(full_tags)
    finally:
        if not args.keep_binary and binary_path.exists():
            try:
                binary_path.unlink()
            except OSError:
                pass

    return 0


if __name__ == "__main__":  # pragma: no cover - CLI entry point
    try:
        raise SystemExit(main())
    except RuntimeError as exc:
        print(f"error: {exc}", file=sys.stderr)
        raise SystemExit(1)
