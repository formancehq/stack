<!-- PROJECT LOGO -->
<p align="center">
  <a href="https://github.com/formancehq/ledger">
    <img src="https://www.formance.com/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Fdiagram-10.adb193c2.svg&w=3840&q=75" alt="Formance">
  </a>

  <h1 align="center">Formance Stack</h1>

  <p align="center">
    Open, modular foundation for unique payments flows
    <br />
    <br />
    Build, operate and track money movements of any size
    <br />
    and shape on a solid and open-source system-of-record.
    <br />
    <br />
    <a href="https://www.formance.com/slack">Slack</a>
    ·
    <a href="https://www.formance.com.com">Website</a>
  </p>
</p>

<p align="center">
   <a href="https://github.com/formancehq/ledger/stargazers"><img src="https://img.shields.io/github/stars/formancehq/ledger" alt="Github Stars"></a>
   <a href="https://github.com/formancehq/ledger/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-mit-purple" alt="License"></a>
   <a href="https://www.ycombinator.com/companies/formance-fka-numary"><img src="https://img.shields.io/badge/Backed%20by-Y%20Combinator-%23f26625"></a>
</p>

## 💻 Deploy locally

### Requirements
1. Install Docker on your machine;
2. Make sure Docker Compose is installed and available (it should be the case if you have chosen to install Docker via Docker Desktop); and
3. Make sure Git is installed on your machine.


### Run the app
To start using Formance Stack, run the following commands in a shell:

```
# Get the code
git clone https://github.com/formancehq/stack.git

# Go to Lago folder
cd stack

# Start
docker-compose up -d
```

You can now open your browser and go to http://localhost to connect to the application. Stack's API is exposed at http://localhost/api.
