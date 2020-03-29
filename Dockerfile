FROM thomasweise/docker-texlive-full:latest

EXPOSE 8080

COPY dicesum /dicesum

CMD ["/dicesum"]