#/usr/bin/env bash

# This Script generates a Response Policy Zone file from PiHole Sources
# Response Policy: https://en.wikipedia.org/wiki/Response_policy_zone

# set up lists
TRIM_LIST="https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts "
TRIM_LIST+="http://sysctl.org/cameleon/hosts "
TRIM_LIST+="https://hosts-file.net/ad_servers.txt "
SIMPLE_LIST="https://mirror1.malwaredomains.com/files/justdomains "
SIMPLE_LIST+="https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist "
SIMPLE_LIST+="https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt "
SIMPLE_LIST+="https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt "

# Set up whitelist
WHITELIST="thepiratebay\.org sendgrid\.net googleadservices\.com doubleclick\.net"

echo "Creating Temporary File"
TMP_FILE=$(mktemp /tmp/pihole-hosts.XXXXXX)
BLOCKLIST_FILE="/root/blocklist"
DNS_ZONE_FILE="/var/named/pihole-ads"

echo "Pulling Complex Hosts"
for i in ${TRIM_LIST}; do
	curl --silent ${i} | awk '{print $2}' >> ${TMP_FILE}
	if [[ ${?} -ne 0 ]]; then
		echo "Failed downloading from ${i}"
	fi
done

echo "Pulling Simple Hosts"
for i in ${SIMPLE_LIST}; do
	curl --silent ${i} >> ${TMP_FILE}
	if [[ ${?} -ne 0 ]]; then
		echo "Failed downloading from ${i}"
	fi
done

echo "Blocklist downloaded: $(wc -l ${TMP_FILE} | awk '{print $1}')"

TOTAL_LINES=$(wc -l ${TMP_FILE} | awk '{print $1}')
sort -u ${TMP_FILE} > ${BLOCKLIST_FILE}
//" ${BLOCKLIST_FILE}
echo "$((${TOTAL_LINES} - $(wc -l ${BLOCKLIST_FILE} | awk '{print $1}'))) Duplicates Removed"

TOTAL_LINES=$(wc -l ${TMP_FILE} | awk '{print $1}')
grep -E "^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$" ${BLOCKLIST_FILE} > ${TMP_FILE}
sed -i -E "/[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/d" ${TMP_FILE}
echo "$((${TOTAL_LINES} - $(wc -l ${TMP_FILE} | awk '{print $1}'))) Invalid Lines Removed"

cat ${TMP_FILE} > ${BLOCKLIST_FILE}

rm -f ${TMP_FILE}
echo "Temporary File Removed"

echo "JFRUTCAyNGgKQAlJTglTT0EJbnMxLnRlc3QtY2hhbWJlci0xMy5sYW4uCXNwaWRlci5zbW9vdGhuZXQub3JnLgkoCgkJMjAxODEyMjMwMAk7CVNlcmlhbAoJCTFoCQk7CVJlZnJlc2gKCQkzMG0JCTsJUmV0cnkKCQkxdwkJOwlFeHBpcmUKCQkxaAkJOwlNaW5pbXVtCikKCjsKOwlOYW1lIFNlcnZlcnMKOwoJCQkJCQlJTglOUwluczEudGVzdC1jaGFtYmVyLTEzLmxhbi4KCjsKOyBBZGRyZXNzZXMKOwo=" | base64 -di > ${DNS_ZONE_FILE}
sed -i -E "s/[0-9]+\t;\tSerial/$(date +%Y%m%d%H)\t;\tSerial/" ${DNS_ZONE_FILE}
for i in $(cat ${BLOCKLIST_FILE}); do
	MATCH=0
	for WL in ${WHITELIST}; do
		if [[ ${i} =~ ${WL} ]]; then
			MATCH=1
		fi
	done

	if [[ ${MATCH} -eq 0  ]]; then
		echo "${i}                      IN      CNAME   ns1.test-chamber-13.lan." >> ${DNS_ZONE_FILE}
	fi
done

echo "Restarting named"
systemctl reload named-chroot
