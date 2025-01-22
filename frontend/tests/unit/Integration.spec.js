const { mount } = require('@vue/test-utils')
const { describe, it, expect, beforeEach } = require('@jest/globals')
const CodeEditor = require('@/components/CodeEditor.vue')

describe('Backend Integration Tests', () => {
  let wrapper

  beforeEach(() => {
    wrapper = mount(CodeEditor)
  })

  it('successfully executes Python code', async () => {
    const code = 'print("Hello, World!")'
    
    await wrapper.setData({
      code,
      language: 'python'
    })

    await wrapper.vm.executeCode()

    // Wait for the API call to complete
    await new Promise(resolve => setTimeout(resolve, 100))

    // Check if output contains the expected result
    expect(wrapper.vm.output).toContain('Hello, World!')
    expect(wrapper.vm.error).toBe('')
  })

  it('successfully executes JavaScript code', async () => {
    const code = 'console.log("Hello from JS")'
    
    await wrapper.setData({
      code,
      language: 'js'
    })

    await wrapper.vm.executeCode()

    // Wait for the API call to complete
    await new Promise(resolve => setTimeout(resolve, 100))

    // Check if output contains the expected result
    expect(wrapper.vm.output).toContain('Hello from JS')
    expect(wrapper.vm.error).toBe('')
  })

  it('handles syntax errors gracefully', async () => {
    const code = 'print("Hello' // Missing closing quote
    
    await wrapper.setData({
      code,
      language: 'python'
    })

    await wrapper.vm.executeCode()

    // Wait for the API call to complete
    await new Promise(resolve => setTimeout(resolve, 100))

    // Check if error is displayed
    expect(wrapper.vm.error).toBeTruthy()
    expect(wrapper.vm.output).toBe('')
  })

  it('handles rate limiting', async () => {
    // Make multiple requests in quick succession
    const requests = Array(15).fill().map(async () => {
      await wrapper.setData({
        code: 'print("test")',
        language: 'python'
      })
      return wrapper.vm.executeCode()
    })

    await Promise.all(requests)

    // Wait for all requests to complete
    await new Promise(resolve => setTimeout(resolve, 100))

    // Check if rate limit error is received
    const hasRateLimitError = wrapper.vm.error.includes('rate limit') || 
                             wrapper.vm.error.includes('too many requests')
    expect(hasRateLimitError).toBe(true)
  })

  it('handles large code input', async () => {
    // Generate a large code string
    const largeCode = 'print("x")\n'.repeat(2000)
    
    await wrapper.setData({
      code: largeCode,
      language: 'python'
    })

    await wrapper.vm.executeCode()

    // Wait for the API call to complete
    await new Promise(resolve => setTimeout(resolve, 100))

    // Check if error about code size is displayed
    expect(wrapper.vm.error).toContain('length exceeds')
    expect(wrapper.vm.output).toBe('')
  })

  it('blocks dangerous system commands', async () => {
    const dangerousCode = 'import os; os.system("rm -rf /")'
    
    await wrapper.setData({
      code: dangerousCode,
      language: 'python'
    })

    await wrapper.vm.executeCode()

    // Wait for the API call to complete
    await new Promise(resolve => setTimeout(resolve, 100))

    // Check if security error is displayed
    expect(wrapper.vm.error).toContain('system')
    expect(wrapper.vm.output).toBe('')
  })
}) 