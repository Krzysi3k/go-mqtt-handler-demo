for f in $(find . -type f -name "*.go"); do
    sed -r -i 's/[a-fA-F0-9]{2}([:-][a-fA-F0-9]{2}){5}/macaddr/;s/([0-9]{1,3}\.){3}[0-9]{1,3}/ipaddr/' "$f"
done