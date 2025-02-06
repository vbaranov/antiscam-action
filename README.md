# antiscam-action

GitHub action to automatically hide scam comments on issues.

## Usage

Add the following to `.github/workflows/antiscam.yml` in your repository:

```
name: antiscam

on:
  issue_comment:
    types:
      - created
      - edited

permissions:
  pull-requests: write
  issues: write

jobs:
  build:
    name: Antiscam
    runs-on: ubuntu-latest

    steps:
      - uses: vbaranov/antiscam-action@main
        with:
          token: ${{ github.token }}
        env:
          SCAM_ACTION_WHITELISTED_LOGINS: ${{ vars.SCAM_ACTION_WHITELISTED_LOGINS }}
```

`SCAM_ACTION_WHITELISTED_LOGINS` is the Github Actions env variable, which contains comma-separated list of GitHub logins which are whitelisted from antiscam check.

