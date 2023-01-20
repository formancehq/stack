FROM gitpod/workspace-full:2023-01-16-03-31-28

RUN brew install kubectl
RUN curl -sL https://get.garden.io/install.sh | bash
ENV PATH=$PATH:$HOME/.garden/bin
RUN sudo apt install rsync -y