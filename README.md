# ossfrontend

nats reply -s nats://demo.nats.io ns.oss.router --queue worker_group_create "service instance create reply #{{Count}}"