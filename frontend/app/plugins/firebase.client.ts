import { initializeApp } from 'firebase/app'
import { getAuth } from 'firebase/auth'

const firebaseConfig = {
  apiKey: 'AIzaSyClPzXIfMbZVg7lCP5MoIfFlYuV6-Xv1As',
  authDomain: 'cinema-25e75.firebaseapp.com',
  projectId: 'cinema-25e75',
  storageBucket: 'cinema-25e75.firebasestorage.app',
  messagingSenderId: '887880820646',
  appId: '1:887880820646:web:7c8642d1a3581f43d2d167',
}

export default defineNuxtPlugin(() => {
  const app = initializeApp(firebaseConfig)
  const auth = getAuth(app)

  return {
    provide: {
      firebaseAuth: auth,
    },
  }
})
