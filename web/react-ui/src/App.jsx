import { useState } from 'react'
import { Layout, Typography, Input, Select, Button, Card, Alert, Divider, theme } from 'antd'
import { DownloadOutlined, FileZipOutlined, GithubOutlined, GlobalOutlined } from '@ant-design/icons'
import './App.css'

const { Header, Content, Footer } = Layout
const { Title, Paragraph, Link } = Typography
const { Option } = Select

// Encoding options
const encodingOptions = [
  'utf-8', 'gbk', 'gb18030', 'big5', 'euc-jp', 'iso-2022-jp', 'shift-jis',
  'euc-kr', 'utf-16be', 'utf-16le', 'koi8-r', 'koi8-u', 'cp437', 'ibm866',
  'macintosh', 'iso-8859-2', 'iso-8859-3', 'iso-8859-4', 'iso-8859-5',
  'iso-8859-6', 'iso-8859-7', 'iso-8859-8', 'iso-8859-10', 'iso-8859-13',
  'iso-8859-14', 'iso-8859-15', 'iso-8859-16', 'windows-874', 'windows-1250',
  'windows-1251', 'windows-1252', 'windows-1253', 'windows-1254', 'windows-1255',
  'windows-1256', 'windows-1257', 'windows-1258', 'x-mac-cyrillic', 'x-user-defined'
]

function App() {
  const [url, setUrl] = useState('')
  const [charset, setCharset] = useState('utf-8')
  const [files, setFiles] = useState([])
  const [selectedFiles, setSelectedFiles] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  
  const { token } = theme.useToken()

  const handleRead = async () => {
    setError(null)
    setLoading(true)
    
    try {
      new URL(url) // Validate URL
      const queryParams = new URLSearchParams()
      queryParams.set('charset', charset)
      queryParams.set('url', url)
      
      const response = await fetch(`/list?${queryParams.toString()}`)
      
      if (response.status === 500) {
        const message = await response.text()
        setError(message)
      } else {
        const data = await response.json()
        setFiles(data.Files || [])
        setSelectedFiles([])
      }
    } catch (err) {
      setError('Invalid URL or network error: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleDownload = async () => {
    if (selectedFiles.length === 0) {
      setError('Please select at least one file to download')
      return
    }
    
    setError(null)
    setLoading(true)
    
    try {
      const queryParams = new URLSearchParams()
      queryParams.set('charset', charset)
      queryParams.set('url', url)
      
      const response = await fetch(`/pack?${queryParams.toString()}`, {
        method: 'POST',
        body: JSON.stringify(selectedFiles),
      })
      
      if (response.status === 500) {
        const message = await response.text()
        setError(message)
      } else {
        const fileStream = window.streamSaver.createWriteStream('package.zip')
        const readableStream = response.body
        
        if (window.WritableStream && readableStream.pipeTo) {
          await readableStream.pipeTo(fileStream)
        } else {
          const writer = fileStream.getWriter()
          const reader = response.body.getReader()
          
          const pump = async () => {
            const { done, value } = await reader.read()
            if (done) {
              writer.close()
            } else {
              await writer.write(value)
              await pump()
            }
          }
          
          await pump()
        }
      }
    } catch (err) {
      setError('Download error: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleUrlChange = (e) => {
    setUrl(e.target.value)
    setFiles([])
    setSelectedFiles([])
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{
        display: 'flex',
        alignItems: 'center',
        background: token.colorBgContainer,
        borderBottom: `1px solid ${token.colorBorder}`,
        padding: '0 50px'
      }}>
        <FileZipOutlined style={{ fontSize: '32px', color: token.colorPrimary, marginRight: '16px' }} />
        <Title level={2} style={{ margin: 0, color: token.colorPrimary }}>
          Archive Proxy
        </Title>
      </Header>
      
      <Content style={{ padding: '50px' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          <Card
            title={
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <GlobalOutlined style={{ marginRight: '8px' }} />
                Read Remote Archive Files
              </div>
            }
            style={{ marginBottom: '24px' }}
            extra={
              <Link href="https://github.com/Heng-Bian/archive-proxy" target="_blank">
                <GithubOutlined style={{ fontSize: '20px' }} />
              </Link>
            }
          >
            <Paragraph>
              An archive proxy written in Go language supporting zip, tar, 7z, rar (including rar5).
            </Paragraph>
            <Paragraph>
              For more information visit{' '}
              <Link href="https://github.com/Heng-Bian/archive-proxy/blob/main/README.md" target="_blank">
                GitHub Repository
              </Link>
            </Paragraph>
          </Card>

          <Card style={{ marginBottom: '24px' }}>
            <div style={{ marginBottom: '24px' }}>
              <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                Archive URL
              </label>
              <Input
                placeholder="Enter archive URL (e.g., https://example.com/file.zip)"
                value={url}
                onChange={handleUrlChange}
                size="large"
                prefix={<GlobalOutlined />}
              />
            </div>

            <div style={{ marginBottom: '24px' }}>
              <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                Character Encoding
              </label>
              <Select
                value={charset}
                onChange={setCharset}
                style={{ width: '100%' }}
                size="large"
              >
                {encodingOptions.map(enc => (
                  <Option key={enc} value={enc}>{enc}</Option>
                ))}
              </Select>
            </div>

            <Button
              type="primary"
              size="large"
              icon={<FileZipOutlined />}
              onClick={handleRead}
              loading={loading}
              block
            >
              List Archive Files
            </Button>
          </Card>

          {error && (
            <Alert
              message="Error"
              description={error}
              type="error"
              closable
              onClose={() => setError(null)}
              style={{ marginBottom: '24px' }}
            />
          )}

          {files.length > 0 && (
            <Card
              title={
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <FileZipOutlined style={{ marginRight: '8px' }} />
                  Archive Contents ({files.length} items)
                </div>
              }
              style={{ marginBottom: '24px' }}
            >
              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                  Select files to download (use Ctrl/Cmd + Click for multiple selection)
                </label>
                <Select
                  mode="multiple"
                  placeholder="Select files to download"
                  value={selectedFiles}
                  onChange={setSelectedFiles}
                  style={{ width: '100%' }}
                  size="large"
                  maxTagCount="responsive"
                >
                  {files.map(file => (
                    <Option key={file} value={file}>
                      {file}
                    </Option>
                  ))}
                </Select>
              </div>

              <Divider />

              <Button
                type="primary"
                size="large"
                icon={<DownloadOutlined />}
                onClick={handleDownload}
                loading={loading}
                disabled={selectedFiles.length === 0}
                block
              >
                Download Selected Files ({selectedFiles.length} selected)
              </Button>
            </Card>
          )}
        </div>
      </Content>
      
      <Footer style={{ textAlign: 'center', background: token.colorBgContainer }}>
        <Paragraph style={{ margin: 0 }}>
          Archive Proxy © 2026 | Built with React & Ant Design
        </Paragraph>
      </Footer>
    </Layout>
  )
}

export default App
