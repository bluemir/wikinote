## 기능 목록

- Wiki 조회
	- markdown rendering
	- 이미지
	- Video
- Wiki 수정
	- 미리보기
	- 이미지 업로드
- Wiki 검색
- Plugin
	- Viewer
		- footer
		- 댓글 Plugin(https://giscus.app/)
		- TAG
		- Recently updated
		- git plugin
- 권한 시스템
	- 아무나 읽되, 쓰는 사람은 일부
	- login
	- 다른 ID 로 login
	- 가입
	- 가입시 Role 지정(config)
	- RBAC

## 기타 고려사항
- cookie free
	- basic auth

## TODO
- [ ] UI 구조개선
	- CustomElement(not React)
	- muti page application
		- non-client routing
- [x] 파일 업로드
- [x] image viewer
- [x] image uploader
- [ ] file upload with drag and drop
- [x] Edit 에서 tab key 수정
- [x] image viewer max width
- [x] video viewer
- [x] Profile
- [ ] Logout(hidden from menu, for debug)
- [x] Unauthroized page
- [ ] Register Page
- [ ] Forbidden page...
- [ ] git plugin
- [ ] html templates 의 not found Error handling
- [ ] http2
- [ ] auto TLS
- [ ] cli options for http/https
	- validation
