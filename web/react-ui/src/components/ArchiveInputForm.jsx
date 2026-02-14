import { Card, Input, Select, Button } from 'antd'
import { GlobalOutlined, FileZipOutlined } from '@ant-design/icons'

function ArchiveInputForm({ url, charset, loading, onUrlChange, onCharsetChange, onRead, encodingOptions }) {
  return (
    <Card style={{ marginBottom: '24px' }}>
      <div style={{ marginBottom: '24px' }}>
        <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
          Archive URL
        </label>
        <Input
          placeholder="Enter archive URL (e.g., https://example.com/file.zip)"
          value={url}
          onChange={onUrlChange}
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
          onChange={onCharsetChange}
          style={{ width: '100%' }}
          size="large"
          options={encodingOptions}
        />
      </div>

      <Button
        type="primary"
        size="large"
        icon={<FileZipOutlined />}
        onClick={onRead}
        loading={loading}
        block
      >
        List Archive Files
      </Button>
    </Card>
  )
}

export default ArchiveInputForm
