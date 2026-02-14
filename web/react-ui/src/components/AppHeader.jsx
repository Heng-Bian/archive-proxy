import { Layout, Typography, theme } from 'antd'
import { FileZipOutlined } from '@ant-design/icons'

const { Header } = Layout
const { Title } = Typography

function AppHeader() {
  const { token } = theme.useToken()

  return (
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
  )
}

export default AppHeader
