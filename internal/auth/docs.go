package auth

/*

roles:
- name: guest
  rules:
  - resources:
      kind: md
      draft: !true
    verbs:
    - read
  - resources:
      kind: history
    verbs:
    - read
  - verbs:
    - search
- name: admin
  rules:
  - verbs:
    - create
	- read
	- update
	- delete
	- search

*/
