<template>
    <fieldset class='login-form'>
        <legend>로그인</legend>
        <div class='input-wrap'>
            <input type="text" id="email" v-model="email" placeholder="아이디">
            <input type="password" id="password" v-model="password" placeholder="비밀번호">
        </div>
        <div class='button-wrap'>
            <button class='submit-btn' v-on:click="submit">로그인</button>
        </div>
    </fieldset>
</template>

<script>
export default {
    name: 'LoginInput',
    data() {
        return {
            email : '',
            password : '',
        }
    },
    methods: {
        submit() {
            console.log(this.email, this.password);
            console.log(this.signin);

        }
    },
    computed: {
        signin(){
            if(this.email == "qjadn0914@naver.com" && this.password == "1q2w3e4r"){
                return {
                    "result": {
                        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjUsImVtYWlsIjoicWphZG4wOTE0QG5hdmVyLmNvbSIsImFkbWluIjpmYWxzZSwiZXhwIjoxNTg1NjM3NTE0fQ.6YsgNJLe3NyCNg28bSKuP66EIlFRMszwEG5GcroolIM"
                    },
                    "success": true,
                    "error": {}
                }
            }
        }
    },
}
</script>

<style scoped>
    .login-form {

    }
    .input-wrap {

    }
    .input-wrap #email {

    }
    .input-wrap #password {

    }
    .button-wrap {

    }
    .button-wrap .submit-btn {
        
    }
    input { display: block; }
</style>

