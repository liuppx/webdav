/**
 * WebDAV Web3 Client Library
 */
class WebDAVWeb3Client {
    constructor(apiBase = 'http://localhost:6065') {
        this.apiBase = apiBase;
        this.token = null;
        this.address = null;
        this.provider = null;
        this.signer = null;
    }

    /**
     * 连接钱包
     */
    async connectWallet() {
        if (typeof window.ethereum === 'undefined') {
            throw new Error('MetaMask is not installed');
        }

        this.provider = new ethers.providers.Web3Provider(window.ethereum);
        await this.provider.send("eth_requestAccounts", []);
        this.signer = this.provider.getSigner();
        this.address = await this.signer.getAddress();

        return this.address;
    }

    /**
     * 认证
     */
    async authenticate() {
        if (!this.address) {
            throw new Error('Wallet not connected');
        }

        // 1. 获取挑战
        const challenge = await this.getChallenge(this.address);

        // 2. 签名
        const signature = await this.signer.signMessage(challenge.message);

        // 3. 验证
        const result = await this.verifySignature(this.address, signature);

        this.token = result.token;

        return result;
    }

    /**
     * 获取挑战
     */
    async getChallenge(address) {
        const response = await fetch(
            `${this.apiBase}/api/auth/challenge?address=${address}`
        );

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to get challenge');
        }

        return await response.json();
    }

    /**
     * 验证签名
     */
    async verifySignature(address, signature) {
        const response = await fetch(`${this.apiBase}/api/auth/verify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                address: address,
                signature: signature,
            }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Verification failed');
        }

        return await response.json();
    }

    /**
     * 列出目录
     */
    async listDirectory(path = '/') {
        return await this.request('PROPFIND', path, null, {
            'Depth': '1',
        });
    }

    /**
     * 上传文件
     */
    async uploadFile(path, content) {
        return await this.request('PUT', path, content);
    }

    /**
     * 下载文件
     */
    async downloadFile(path) {
        const response = await this.request('GET', path);
        return await response.text();
    }

    /**
     * 删除文件
     */
    async deleteFile(path) {
        return await this.request('DELETE', path);
    }

    /**
     * 创建目录
     */
    async createDirectory(path) {
        return await this.request('MKCOL', path);
    }

    /**
     * 发送请求
     */
    async request(method, path, body = null, extraHeaders = {}) {
        if (!this.token) {
            throw new Error('Not authenticated');
        }

        const headers = {
            'Authorization': `Bearer ${this.token}`,
            ...extraHeaders,
        };

        const response = await fetch(`${this.apiBase}${path}`, {
            method: method,
            headers: headers,
            body: body,
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        return response;
    }
}

// 使用示例
async function example() {
    const client = new WebDAVWeb3Client();

    try {
        // 1. 连接钱包
        const address = await client.connectWallet();
        console.log('Connected:', address);

        // 2. 认证
        const auth = await client.authenticate();
        console.log('Authenticated:', auth.user);

        // 3. 列出目录
        const files = await client.listDirectory('/');
        console.log('Files:', files);

        // 4. 上传文件
        await client.uploadFile('/test.txt', 'Hello, Web3!');
        console.log('File uploaded');

        // 5. 下载文件
        const content = await client.downloadFile('/test.txt');
        console.log('File content:', content);

    } catch (error) {
        console.error('Error:', error);
    }
}

