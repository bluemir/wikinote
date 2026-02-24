## TODO


- [ ] markdown에서 math 표현 rendering
- [ ] router 에서 api 들은 api() 로 감싸기
	- 0xC0DE 참조
- [ ] test 추가
	- 특히 markdown rendering
		- footnote
		- table
		- checklist
- [ ] markdown css 손보기
- [ ] layout 용 custom element 는 제거
	- layout은 html 로 하고, component 들만 custom element로
- [x] gnb는 position: sticky 로
- [x] text-area `field-size: content` 적용
- [ ] Login page, Register 페이지 업데이트
- [x] pubsub을 0xC0DE에 있는 버전으로 업데이트
- [x] `github.com/cockroachdb/errors` 로 교체
- [x] embed 파일을 assets 내로 이동
- [ ] 현재 `d1fc296` 를 서빙 중인데 차이점 및 버전업 할때 고려사항 확인
- [ ] password 암호화 방식 변경.
- [ ] web.mk 의 정리
	- embed.go 에서 build 하므로 해당 script 는 필요하면 정리.
- [x] functional package 를 github.com/bluemir/funtional 로 대체
