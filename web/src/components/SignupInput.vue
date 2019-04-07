<template>
    <fieldset class='signup-form'>
        <legend>회원가입</legend>
        <div class='input-wrap'>
            <input type="text" id="email" v-on:keyup.enter="enter" v-on:blur="blur" v-model="email" placeholder="아이디" maxlength="50">
            <input type="password" id="password" v-on:keyup.enter="enter" v-model="password" placeholder="비밀번호" maxlength="50">
            <input type="password" id="passconfirm" v-on:keyup.enter="enter" v-model="passconfirm" placeholder="비밀번호확인" maxlength="50">
        </div>
        <div class='button-wrap'>
            <button class='submit-btn' v-on:click="submit">회원가입</button>
        </div>
    </fieldset>
</template>

<script>
export default {
    name: 'SignupInput',
    data() {
        return {
            email : '',
            password : '',
            passconfirm : '',
        }   
    },
    methods: {        
        enter(e){
            console.log(e.target);
            if(e.target.id == "email"){
                document.querySelector("#password").focus();
            } else if(e.target.id == "password"){
                document.querySelector("#passconfirm").focus();
            } else if(e.target.id == "passconfirm"){
                this.submit();
                    
            }
        },
        blur(e){
            console.log(e);
        },
        chkInput(){
            console.log(this);
            const emailRegEx = /^[0-9a-zA-Z_-]+(\.[0-9a-zA-Z_-]+)*@([0-9a-zA-Z_-]+)(\.[0-9a-zA-Z_-]+){1,2}$/;
            const chkEmail = emailRegEx.test(this.email);
            let flag = true;
            if(!chkEmail){
                alert("이메일 패턴이 맞지 않습니다.");
                document.querySelector("#email").focus();
                flag = false;
            } else if(this.password.length <= 0){
                //비밀번호 크기 (몇자 이상 몇자 미만)
                //비밀번호 조합 고민해야할듯
                document.querySelector("#password").focus();
                alert("비밀번호를 입력해주세요.");
                flag = false;
            }

            return flag;
        },
        submit() {
            console.log(this.chkInput());
            return;
            const params = {
                email: this.email,
                password: this.password
            }

            let url = '';
            // url = '192.168.1.3:8080';
            // url = 'www.sumwhere.kr';
            axios.post(url+'/v1/api/signin',params)
            .then(res => {
                console.log(res);
                if(res.success){
                    const token = res.result.token;

                    //쿠키만료일 설정 => 만료일 = 7일;
                    const date = new Date();
                    date.setDate(date.getDate() + 7);
                    //

                    document.cookie = 'jwt='+escape(token)+'; expires='+ date.toUTCString();
                    console.log('cookie', document.cookie);
                } else {
                    alert(this.error);
                }
            })
            .catch(err => {
                console.error(err); 
            })
            
        }
    },
}
</script>

<style scoped>
    .login-form {
        width: 460px;
        height: 270px;
        margin: 0 auto;
    }
    .input-wrap input {
        width: calc(100% - 24px);
        height: 30px;
        padding: 10px;
        margin-bottom: 15px;
    }
    .button-wrap {

    }
    .button-wrap .submit-btn {
        width: 100%;
        height: 60px;
        margin: 15px 0;

        background-color: dodgerblue;
        border: 1px solid dodgerblue;

        color: #ffffff;
        font-size: 20px;
    }
    input { display: block; }
</style>

