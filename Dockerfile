FROM thomasweise/docker-texlive-full:latest

EXPOSE 8080

COPY dicemath /dicemath

CMD ["/dicemath"]