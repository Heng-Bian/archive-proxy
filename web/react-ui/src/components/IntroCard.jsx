import { Card, Typography } from 'antd'
import { GlobalOutlined, GithubOutlined } from '@ant-design/icons'

const { Paragraph, Link } = Typography

function IntroCard() {
  return (
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
  )
}

export default IntroCard
