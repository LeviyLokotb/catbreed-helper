#!/bin/bash
wget https://github.com/microsoft/onnxruntime/releases/download/v1.25.0/onnxruntime-linux-x64-1.25.0.tgz
tar -xzf onnxruntime-linux-x64-1.25.0.tgz

cp onnxruntime-linux-x64-1.25.0/lib/libonnxruntime* /usr/lib/

ln -sf /usr/lib/libonnxruntime.so.1.25.0 /usr/lib/onnxruntime.so

ldconfig