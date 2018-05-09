FROM scratch
ADD main /
ADD static /static
ADD frontend/templates /frontend/templates
CMD ["/main"]