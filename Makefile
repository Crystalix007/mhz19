all: co2 co2-logger-prometheus

co2: cmd/mhz19/main.go
	go build -o $@ $<

co2-logger-prometheus: cmd/mhz19/prometheus/main.go
	go build -o $@ $<

PREFIX ?= /usr

.PHONY: install
install: all
	install -Dm755 co2-logger-prometheus $(DESTDIR)$(PREFIX)/bin/
	install -Dm755 co2 $(DESTDIR)$(PREFIX)/bin/
	install -Dm644 co2-prometheus.service $(DESTDIR)/etc/systemd/system/
