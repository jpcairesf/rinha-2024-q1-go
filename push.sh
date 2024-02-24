docker build -t jpcairesf/rinha-2024-q1-go:latest . \
  && docker push jpcairesf/rinha-2024-q1-go:latest \
  && docker-compose down --volumes \
  && docker-compose up