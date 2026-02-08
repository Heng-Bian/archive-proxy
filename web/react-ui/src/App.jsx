import { useState } from 'react'
import { Layout, Typography, Input, Select, Button, Card, Alert, Divider, Tree, theme } from 'antd'
import { DownloadOutlined, FileZipOutlined, GithubOutlined, GlobalOutlined, FolderOutlined, FileOutlined } from '@ant-design/icons'
import './App.css'

const { Header, Content, Footer } = Layout
const { Title, Paragraph, Link } = Typography

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

// Transform encoding options for Select component
const encodingSelectOptions = encodingOptions.map(enc => ({ label: enc, value: enc }))

// Convert flat file list to tree structure
function buildTreeData(files) {
  const root = { children: {} }
  
  files.forEach(filePath => {
    const isDir = filePath.endsWith('/')
    // Filter out empty strings from split to handle edge cases like '//' or trailing slashes
    const parts = filePath.split('/').filter(p => p !== '')
    
    let current = root
    parts.forEach((part, index) => {
      const isLastPart = index === parts.length - 1
      
      // Use consistent key (always without trailing slash for the part name)
      if (!current.children[part]) {
        current.children[part] = {
          key: parts.slice(0, index + 1).join('/') + (isLastPart && isDir ? '/' : ''),
          title: part,
          isLeaf: isLastPart && !isDir,
          children: {}
        }
      } else if (isLastPart && isDir) {
        // Update the key to include trailing slash if this is a directory entry
        current.children[part].key = parts.slice(0, index + 1).join('/') + '/'
        current.children[part].isLeaf = false
      }
      
      current = current.children[part]
    })
  })
  
  // Convert children object to array recursively
  function convertToArray(node) {
    const childrenArray = Object.values(node.children).map(child => {
      const converted = {
        key: child.key,
        title: child.title,
        icon: child.isLeaf ? <FileOutlined /> : <FolderOutlined />,
        isLeaf: child.isLeaf
      }
      
      if (!child.isLeaf && Object.keys(child.children).length > 0) {
        converted.children = convertToArray(child)
      }
      
      return converted
    })
    
    // Sort: folders first, then files, both alphabetically
    return childrenArray.sort((a, b) => {
      if (a.isLeaf === b.isLeaf) {
        return a.title.localeCompare(b.title)
      }
      return a.isLeaf ? 1 : -1
    })
  }
  
  return convertToArray(root)
}

function App() {
  const [url, setUrl] = useState('')
  const [charset, setCharset] = useState('utf-8')
  const [files, setFiles] = useState([])
  const [treeData, setTreeData] = useState([])
  const [selectedFiles, setSelectedFiles] = useState([])
  const [checkedKeys, setCheckedKeys] = useState([])
  const [expandedKeys, setExpandedKeys] = useState([])
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
        const fileList = data.Files || []
        setFiles(fileList)
        setTreeData(buildTreeData(fileList))
        setSelectedFiles([])
        setCheckedKeys([])
        setExpandedKeys([])
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
    setTreeData([])
    setSelectedFiles([])
    setCheckedKeys([])
    setExpandedKeys([])
  }

  const onCheck = (checkedKeysValue) => {
    // Filter out directory keys (ending with /)
    const fileKeys = checkedKeysValue.filter(key => !key.endsWith('/'))
    setCheckedKeys(checkedKeysValue)
    setSelectedFiles(fileKeys)
  }

  const onExpand = (expandedKeysValue) => {
    setExpandedKeys(expandedKeysValue)
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
                options={encodingSelectOptions}
              />
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
                  Select files to download (directories shown for structure only)
                </label>
                <div style={{ 
                  border: `1px solid ${token.colorBorder}`,
                  borderRadius: token.borderRadius,
                  padding: '16px',
                  maxHeight: '500px',
                  overflow: 'auto',
                  backgroundColor: token.colorBgContainer
                }}>
                  <Tree
                    checkable
                    selectable={false}
                    checkedKeys={checkedKeys}
                    expandedKeys={expandedKeys}
                    onCheck={onCheck}
                    onExpand={onExpand}
                    treeData={treeData}
                    showIcon
                  />
                </div>
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
