package auth

/*
kind: role
name: guest
rules:
- objects:
  - "kind==md && draft!=true"
  - "kind==history"
  actions:
  - app:read
- actions:
  - app:search

---
kind: user
name: bluemir
labels:
  role/guest: true
  role/admin: true

---
kind: default-role
roles:
- guest
*/
