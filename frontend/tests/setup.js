const { jest } = require('@jest/globals')
const { config } = require('@vue/test-utils')

// Mock Monaco Editor
jest.mock('@guolao/vue-monaco-editor', () => ({
  default: {
    name: 'MonacoEditor',
    template: '<div class="monaco-editor"></div>'
  }
}))

// Mock Firebase
jest.mock('firebase/app', () => ({
  initializeApp: jest.fn(),
  getApps: jest.fn(() => []),
  getApp: jest.fn()
}))

jest.mock('firebase/auth', () => ({
  getAuth: jest.fn(),
  signInWithPopup: jest.fn(),
  GoogleAuthProvider: jest.fn()
}))

// Configure Vue Test Utils
config.global.mocks = {
  $store: {
    state: {},
    commit: jest.fn(),
    dispatch: jest.fn()
  }
}

// Mock fetch globally
global.fetch = jest.fn()

// Mock console methods to avoid noise in tests
console.error = jest.fn()
console.warn = jest.fn() 