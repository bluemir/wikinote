# 권한 관리

## 용어정리

* 기능 : API 또는 Page
* 기능의 모음 : Action

하나의 API는 여러개의 Action에 속할수 있고
하나의 Action 은 여러개의 API로 이루어 진다.

이 기능의 모음은 보통 code 작성 당시에 알수 있으므로 code 로 define 한다.

* 사용자 : User
* 사용자의 모음 : Role

사용자는 하나의 Role에만 할당 될수 있다.
대신 Role은 동적으로 늘릴수 있다. 만약 두가지 Role을 가지고 싶다면, 새로운 Role을 추가해야한다.(겸직 불가 ^^)
이렇게 함으로써 유연하고 심플한 관리가 된다.

Role에 어떤 Action이 *허가* 되는 것인지를 Rule 이라고 한다.

부연 설명:
Action과 기능이 n:m의 관계이므로 허가된 동작과 허가 되지 않는 동작이 겹칠수 있다.
이 경우에 허가가 우선한다.
왜냐하면 허가하는 것을 체크 하기 때문에 어느 하나에서라도 허가 되면 허가 된것으로 취급하기 때문이다.

Rule은 제한 없이 추가가 가능하다.

## Default Rule

기본 Role은 다음과 같다.
* Admin
* Editor
* Reader
* Guest

기본설정상 처음 가입한 사용자는 Reader에 속한다.

기본 Rule 은 다음과 같이 정의 된다.
{
	Role: Admin,
	Allow: view, search, edit, attach, user,

	Role: Editor,
	Allow: view, search, edit, attach

	Role: Reader,
	Allow: view, search

	Role: Guest,
	Allow: view
}
