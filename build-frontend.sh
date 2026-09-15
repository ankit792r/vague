#!/bin/bash

set -e

cd frontend
bun install
bun run build
cd ..
