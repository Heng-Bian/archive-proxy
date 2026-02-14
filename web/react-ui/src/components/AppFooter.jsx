import { Layout, Typography, theme } from 'antd'

const { Footer } = Layout
const { Paragraph } = Typography

function AppFooter() {
  const { token } = theme.useToken()

  return (
    <Footer style={{ textAlign: 'center', background: token.colorBgContainer }}>
      <Paragraph style={{ margin: 0 }}>
        Archive Proxy © 2026 | Built with React & Ant Design
      </Paragraph>
    </Footer>
  )
}

export default AppFooter
