FROM node:alpine
WORKDIR /app
COPY package*.json .
RUN npm install
COPY . .
RUN npm run build 
EXPOSE 3000
CMD ["npm", "start"]
#docker build --pull --rm -f "Dockerfile" -t serchiboi/frontend_vmnode_p2so1:latest "./"
#docker push serchiboi/frontend_vmnode_p2so1