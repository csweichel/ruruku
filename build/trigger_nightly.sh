#!/bin/bash

git tag --delete nightly || echo "Don't have a nightly tag yet"
git tag -a nightly -m 'Nightly build'

# install deployment key (inspired by https://stackoverflow.com/questions/18935539/authenticate-with-github-using-token)
openssl aes-256-cbc -pbkdf2 -k "$NIGHTLY_KEY_PASSWORD" -d -a -in build/travisNightlyKey.enc -out build/deployKey
chmod 600 build/deployKey
echo -e "Host github.com\n  IdentityFile $PWD/build/deployKey" > ~/.ssh/config
echo "github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==" > ~/.ssh/known_hosts

# swap origin
git remote rm origin
git remote add origin git@github.com:32leaves/ruruku.git

git push origin :nightly
git push origin nightly

echo "Created nightly tag on $(git rev-parse HEAD)"
