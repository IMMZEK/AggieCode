<template>
    <div class="login-container">
      <button @click="signInWithGoogle" class="google-login-button">
        Sign In with Google
      </button>
    </div>
  </template>
  
  <script>
  import { auth } from "../firebase";
  import { signInWithPopup, GoogleAuthProvider } from "firebase/auth";
  import { useRouter } from "vue-router";
  
  export default {
    setup() {
      const router = useRouter();
  
      const signInWithGoogle = async () => {
        const provider = new GoogleAuthProvider();
        try {
          const result = await signInWithPopup(auth, provider);
          // The signed-in user info.
          const user = result.user;
          console.log("User:", user);
  
          // Get the ID token
          const idToken = await user.getIdToken();
          console.log("ID Token:", idToken);
  
          // Redirect to the main app
          router.push("/");
        } catch (error) {
          console.error("Error signing in with Google:", error);
        }
      };
  
      return { signInWithGoogle };
    },
  };
  </script>
  
  <style scoped>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }
  
  .google-login-button {
    padding: 10px 20px;
    background-color: #4285f4;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }
  </style>