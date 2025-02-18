#!/bin/bash

# Install JupyterLab Night theme
pip install --no-cache-dir jupyterlab_night

# Create necessary directories
mkdir -p /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/apputils-extension

# Create theme settings file
cat << EOF > /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/apputils-extension/themes.jupyterlab-settings
{
    "theme": "JupyterLab Night",
    "theme-scrollbars": true
}
EOF

# Create editor settings
mkdir -p /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/fileeditor-extension
cat << EOF > /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/fileeditor-extension/plugin.jupyterlab-settings
{
    "editorConfig": {
        "lineNumbers": true,
        "lineWrap": "on",
        "fontSize": 14,
        "theme": "dark"
    }
}
EOF

# Set permissions
chown -R jovyan:users /home/jovyan/.jupyter