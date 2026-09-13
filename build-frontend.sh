#!/bin/bash

set -e

cd platform/frontend
bun install
bun run build
cd ../..