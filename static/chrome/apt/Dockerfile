ARG BASE_IMAGE=browsers/base:7.4.2
FROM $BASE_IMAGE

ARG VERSION
ARG PACKAGE=google-chrome-stable
ARG INSTALL_DIR=chrome

LABEL browser=$PACKAGE:$VERSION

RUN \
        curl -fsSL https://dl-ssl.google.com/linux/linux_signing_key.pub | gpg --dearmor -o /usr/share/keyrings/google-chrome.gpg && \
        echo 'deb [arch=amd64 signed-by=/usr/share/keyrings/google-chrome.gpg] https://dl.google.com/linux/chrome/deb/ stable main' > /etc/apt/sources.list.d/google.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install iproute2 $PACKAGE=$VERSION && \
         ( \
          sed -i -e 's@exec -a "$0" "$HERE/chrome"@& --no-sandbox --disable-gpu --remote-debugging-port=0@' /opt/google/$INSTALL_DIR/$PACKAGE || \
          sed -i -e 's@exec -a "$0" "$HERE/chrome"@& --no-sandbox --disable-gpu --remote-debugging-port=0@' /opt/google/$INSTALL_DIR/google-chrome \
         ) && \
        chown root:root /opt/google/$INSTALL_DIR/chrome-sandbox && \
        chmod 4755 /opt/google/$INSTALL_DIR/chrome-sandbox && \
        google-chrome --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
