import { Card, Tree, Button, Divider, theme } from 'antd'
import { FileZipOutlined, DownloadOutlined } from '@ant-design/icons'

function FileTreeCard({ 
  files, 
  treeData, 
  checkedKeys, 
  expandedKeys, 
  selectedFiles, 
  loading, 
  onCheck, 
  onExpand, 
  onDownload 
}) {
  const { token } = theme.useToken()

  if (files.length === 0) {
    return null
  }

  return (
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
        onClick={onDownload}
        loading={loading}
        disabled={selectedFiles.length === 0}
        block
      >
        Download Selected Files ({selectedFiles.length} selected)
      </Button>
    </Card>
  )
}

export default FileTreeCard
