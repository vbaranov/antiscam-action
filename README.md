# antiscam-action

GitHub action to automatically hide scam comments on issues.

## Usage

Add the following to `.github/workflows/antiscam.yml` in your repository:

```
name: antiscam

on:
  issues:
    types:
      - opened
      - edited
      - reopened
  pull_request:
    types:
      - opened
      - edited
      - reopened
      - synchronize

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
```

