# Controller deb

## Files

- Controller Binary: `/usr/bin/it-system-controller`
- Configuration: `/etc/it-system/config-controller.yaml`
- System Service: `/lib/systemed/system/it-system-controller.service`
- Frontend Static File: `/usr/share/it-system-controller/frontend`
- DB: `/usr/share/it-system-controller/.db/it-system.db`
- Log: `/usr/share/it-system-controller/.db/log`

## Release Script

```bash
./release_controller.sh <version>
```

This will generate `pkg/it-system-controller-<version>-amd64.deb`.

## Install deb

```bash
sudo dpkg -i it-system-controller-<version>-amd64.deb
```

it-system-controller will be a system service automatically.

If result in dependency issue:

```bash
sudo apt-get -f install

sudo systemctl daemon-reload
sudo systemctl enable --now it-system-controller
sudo systemctl status it-system-controller
```

## Upgrade New deb

```bash
sudo dpkg -i it-system-controller-<version>-amd64.deb
```

## Remove

```bash
sudo dpkg -r it-system-controller
```

This won't remove config and DB files. If you want to remove all the  files:

```bash
sudo dpkg -P it-system-controller
sudo rm -rf /usr/share/it-system-controller/.db
```
