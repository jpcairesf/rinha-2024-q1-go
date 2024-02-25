docker buildx build --platform linux/amd64 -t jpcairesf/rinha-2024-q1-go:latest . \
  && docker push jpcairesf/rinha-2024-q1-go:latest \
  && docker image prune -a --force \
  && docker-compose down --volumes \
  && docker-compose up