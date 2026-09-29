FROM grafana/grafana:13.2.2
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]
COPY assets/logo.svg assets/logo.png /tmp/
USER root
WORKDIR /usr/share/grafana/public/build
RUN sv="$(md5sum /tmp/logo.svg | cut -c1-8)"; pv="$(md5sum /tmp/logo.png | cut -c1-8)"; \
    cp /tmp/logo.svg "static/img/grafana_icon.$sv.svg"; \
    cp /tmp/logo.png "img/fav32.$pv.png"; \
    old="$(ls app.*.js)"; new="app.$sv.js"; \
    sed "s/grafana_icon\.[0-9a-f]*\.svg/grafana_icon.$sv.svg/g" "$old" > "$new"; \
    sed -i "s|$old|$new|g" assets-manifest.json; \
    sed -i "s|href=\"\[\[.FavIcon\]\]\"|href=\"public/build/img/fav32.$pv.png\"|" ../views/index.html
WORKDIR /usr/share/grafana
USER grafana
