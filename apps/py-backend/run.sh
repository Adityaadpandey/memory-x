#!/bin/bash
source venv/bin/activate
uvicorn main:app --host 0.0.0.0 --port 8888
