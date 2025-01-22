const { mount } = require('@vue/test-utils')
const { describe, it, expect, beforeEach, jest } = require('@jest/globals')
const CodeEditor = require('@/components/CodeEditor.vue')

// Mock the code execution API
const mockExecuteCode = jest.fn()
global.fetch = jest.fn()

describe('CodeEditor.vue', () => {
  let wrapper

  beforeEach(() => {
    // Reset all mocks before each test
    jest.clearAllMocks()
    
    // Mock successful API response
    global.fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        output: 'Hello, World!',
        status_message: 'success'
      })
    })

    wrapper = mount(CodeEditor)
  })

  it('renders the code editor', () => {
    expect(wrapper.exists()).toBe(true)
  })

  it('executes code and displays output', async () => {
    // Set the code in the editor
    await wrapper.setData({
      code: 'print("Hello, World!")',
      language: 'python'
    })

    // Trigger code execution
    await wrapper.vm.executeCode()

    // Check if the API was called with correct parameters
    expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/execute', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        language: 'python',
        code: 'print("Hello, World!")',
        method: 'run'
      })
    })

    // Check if output is displayed
    expect(wrapper.vm.output).toBe('Hello, World!')
    expect(wrapper.vm.error).toBe('')
  })

  it('handles execution errors', async () => {
    // Mock API error response
    global.fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        error: 'Syntax error',
        status_message: 'error'
      })
    })

    // Set invalid code
    await wrapper.setData({
      code: 'invalid python code',
      language: 'python'
    })

    // Trigger code execution
    await wrapper.vm.executeCode()

    // Check if error is displayed
    expect(wrapper.vm.error).toBe('Syntax error')
    expect(wrapper.vm.output).toBe('')
  })

  it('handles network errors', async () => {
    // Mock network error
    global.fetch.mockRejectedValue(new Error('Network error'))

    // Set code
    await wrapper.setData({
      code: 'print("Hello")',
      language: 'python'
    })

    // Trigger code execution
    await wrapper.vm.executeCode()

    // Check if error is displayed
    expect(wrapper.vm.error).toContain('Network error')
    expect(wrapper.vm.output).toBe('')
  })

  it('validates language selection', async () => {
    // Set unsupported language
    await wrapper.setData({
      code: 'print("Hello")',
      language: 'invalid'
    })

    // Trigger code execution
    await wrapper.vm.executeCode()

    // Check if error is displayed
    expect(wrapper.vm.error).toContain('Unsupported language')
    expect(wrapper.vm.output).toBe('')
  })
}) 