#!/bin/bash

# Check if configuration already exists
CONFIG_FILE="/home/jovyan/.jupyter/lab/user-settings/@jupyterlab/apputils-extension/themes.jupyterlab-settings"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "First time setup: Installing theme and configuring settings..."
    
    # Install JupyterLab Night theme
    pip install --no-cache-dir jupyterlab_night

    # Create necessary directories
    mkdir -p /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/apputils-extension
    mkdir -p /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/fileeditor-extension

    # Create theme settings file
    cat << EOF > /home/jovyan/.jupyter/lab/user-settings/@jupyterlab/apputils-extension/themes.jupyterlab-settings
{
    "theme": "JupyterLab Night",
    "theme-scrollbars": true
}
EOF

    # Create editor settings
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
    echo "Setting permissions..."
    chown -R jovyan:users /home/jovyan/.jupyter
else
    echo "Configuration already exists, skipping setup..."
fi