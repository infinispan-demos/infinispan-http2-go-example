download:
	wget --show-progress "https://origin-repository.jboss.org/nexus/content/repositories/public-jboss/org/infinispan/server/infinispan-server-build/9.1.0-SNAPSHOT/infinispan-server-build-9.1.0-20170609.185048-78.zip" -O infinispan-server.zip
	unzip -o infinispan-server.zip
	mv infinispan-server-9.1.0-SNAPSHOT infinispan-server
.PHONY: download

prepare-infinispan:
	# We generate keystore for the server
	keytool -genkey -noprompt -trustcacerts -keyalg RSA -alias "localhost" -dname "CN=localhost, OU=Infinispan, O=JBoss, L=Red Hat, ST=World, C=WW" -keypass "secret" -storepass "secret" -keystore "server_keystore.jks"

	# We export the certificate
	keytool -export -keyalg RSA -alias "localhost" -storepass "secret" -file "client_cert.cer" -keystore "server_keystore.jks"

	# Now we need to convert it to PEM
	openssl x509 -inform der -in client_cert.cer -out certificate.pem

	# This will be consumed by Go
	mv server_keystore.jks infinispan-server/standalone/configuration

	# We don't need this
	rm client_cert.cer

	# And the configuration
	cp scripts/standalone-rest-ssl.xml infinispan-server/standalone/configuration
.PHONY: prepare-infinispan

start-infinispan:
	JBOSS_PIDFILE=${PWD}/infinispan.pid LAUNCH_JBOSS_IN_BACKGROUND=true infinispan-server/bin/standalone.sh -c standalone-rest-ssl.xml &
.PHONY: start-infinispan

stop-infinispan:
	cat infinispan.pid | xargs kill
.PHONY: stop-infinispan

run-client:
	go run Main.go
.PHONY: prepare-infinispan