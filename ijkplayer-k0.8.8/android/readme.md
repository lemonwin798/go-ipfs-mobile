
基于ipfs和ijkplayer的播放器示例

在AppActivity.java里使用ipfsApi:

    try {
           
            IpfsApi.api_SetBoostarp("/ip4/122.11.47.95/tcp/4001/ipfs/QmPYpGqwfvcrgfMj6QwcyPS8LC8RJMBBanRu8GDR2bCwCe");
            IpfsApi.api_InitNode(false, "816930a90d14bb804a51f3e7fb867d3f3ce169b688211b5157a0ef40137b167a");
        } catch (Exception e) {
            return;
        }

        try {
            IpfsApi.api_ServeHTTPGateway();
        } catch (Exception e) {
            return;
        }
    }

 
 IpfsApi.api_SetBoostarp("/ip4/122.11.47.95/tcp/4001/ipfs/QmPYpGqwfvcrgfMj6QwcyPS8LC8RJMBBanRu8GDR2bCwCe");
使用私链时，设置引导节点。如果使用ipfs公网，则不需要调用

IpfsApi.api_InitNode(false, "816930a90d14bb804a51f3e7fb867d3f3ce169b688211b5157a0ef40137b167a");
初始化ipfs节点
第一个参数，如果是长期节点，则为false,如果为true则是临时节点
第二个参数，私链的密钥，如果用使用ipfs公网，则填""空字符串即可
