const fs = require('fs');

const requests = [
    {
        method: 'GET',
        url: 'http://localhost:8080/ping',
        example: "pong"
    },
    {
        method: 'POST',
        url: 'http://localhost:8080/lecture',
        data: {
            title: 'Introduction to Computer Science',
            content: 'マークダウン形式 資料へのリンク',
            folderpath: '/School/ComputerScience'
        },
        example: {
            id: 1,
            title: "Introduction to Computer Science",
            content: "マークダウン形式 資料へのリンク",
            folderpath: "/School/ComputerScience"
        }
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/lecture/1',
        example: {
            id: 1,
            title: "Introduction to Computer Science",
            content: "マークダウン形式 資料へのリンク",
            folderpath: "/School/ComputerScience"
        }
    },
    {
        method: 'POST',
        url: 'http://localhost:8080/lecture',
        data: {
            title: 'コンピュータサイエンス入門',
            content: 'content',
            folderpath: '/School/ComputerScience'
        },
        example: {
            id: 2,
            title: 'コンピュータサイエンス入門',
            content: 'content',
            folderpath: '/School/ComputerScience'
        }
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/lecture/byFolder/id/2',
        example: [
            {
                id: 1,
                title: "Introduction to Computer Science",
                content: "マークダウン形式 資料へのリンク",
                folderpath: "/School/ComputerScience"
            },
            {
                id: 2,
                title: 'コンピュータサイエンス入門',
                content: 'content',
                folderpath: '/School/ComputerScience'
            }
        ]
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/lecture/byFolder/path?folderpath=School-ComputerScience',
        example: [
            {
                id: 1,
                title: "Introduction to Computer Science",
                content: "マークダウン形式 資料へのリンク",
                folderpath: "/School/ComputerScience"
            },
            {
                id: 2,
                title: 'コンピュータサイエンス入門',
                content: 'content',
                folderpath: '/School/ComputerScience'
            }
        ]
    },
    {
        method: 'POST',
        url: 'http://localhost:8080/wiki/search',
        data: { query: 'Microsoft', tags: [], resultCount: 10, from: 0 },
        example: [
            {
                "id": 1,
                "type": "sodan",
                "title": "Introduction to Computer Science",
                "Abstract": "This is an introductory course to computer science",
                "createdAt": "2021-01-01 00:00:00",
                "updatedAt": "2021-01-01 00:00:00",
                "ownerTraqId": "kavos",
                "tags": [
                    "ComputerScience"
                ]
            }
        ]
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/wiki/tag?tag=windows',
        example: [
            {
                "id": 1,
                "type": "sodan",
                "title": "Introduction to Computer Science",
                "Abstract": "This is an introductory course to computer science",
                "createdAt": "2021-01-01 00:00:00",
                "updatedAt": "2021-01-01 00:00:00",
                "ownerTraqId": "kavos",
                "tags": [
                    "ComputerScience"
                ]
            }
        ]
    },
    {
        method: 'POST',
        url: 'http://localhost:8080/wiki/tag',
        data: { wikiId: 1, tag: 'ComputerScience' },
        example: {
            "wikiId": 1,
            "tag": "ComputerScience"
        }
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/sodan?wikiId=1',
        example: {
            "id": 1,
            "title": "Introduction to Computer Science",
            "tags": [
                "ComputerScience"
            ],
            "questionMessage": {
                "userTraqId": "kavos",
                "content": "メッセージの中身",
                "createdAt": "2021-01-01 00:00:00",
                "updatedAt": "2021-01-01 00:00:00",
                "stamps": [
                    {
                        "stampId": "abcd-efgh-ijkl",
                        "count": 3
                    }
                ],
                "citations": [
                    {
                        "userTraqId": "kavos",
                        "content": "メッセージの中身",
                        "createdAt": "2021-01-01 00:00:00",
                        "updatedAt": "2021-01-01 00:00:00"
                    }
                ]
            },
            "answerMessages": [
                {
                    "userTraqId": "kavos",
                    "content": "メッセージの中身",
                    "createdAt": "2021-01-01 00:00:00",
                    "updatedAt": "2021-01-01 00:00:00",
                    "stamps": [
                        {
                            "stampId": "abcd-efgh-ijkl",
                            "count": 3
                        }
                    ],
                    "citations": [
                        {
                            "userTraqId": "kavos",
                            "content": "メッセージの中身",
                            "createdAt": "2021-01-01 00:00:00",
                            "updatedAt": "2021-01-01 00:00:00"
                        }
                    ]
                }
            ]
        }
    },
    {
        method: 'POST',
        url: 'http://localhost:8080/memo',
        data: {
            title: 'Introduction to Computer Science',
            content: 'This is an introductory course to computer science',
            tags: ['hoadsoih']
        },
        example: {
            "id": 21,
            "title": "Introduction to Computer Science",
            "content": "This is an introductory course to computer science",
            "ownerTraqId": "kavos",
            "tags": [
                "hoadsoih"
            ],
            "createdAt": "2021-01-01 00:00:00",
            "updatedAt": "2021-01-01 00:00:00"
        }
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/memo/21',
        example: {
            "id": 21,
            "title": "Introduction to Computer Science",
            "content": "This is an introductory course to computer science",
            "ownerTraqId": "kavos",
            "tags": [
                "hoadsoih"
            ],
            "createdAt": "2021-01-01 00:00:00",
            "updatedAt": "2021-01-01 00:00:00"
        }
    },
    {
        method: 'PATCH',
        url: 'http://localhost:8080/memo',
        data: { id: 21, title: 'askdaosjdoa', content: 'This is an introductory course to computer science' },
        example: {
            "id": 21,
            "title": "askdaosjdoa",
            "content": "This is an introductory course to computer science",
            "ownerTraqId": "kavos",
            "tags": [
                "hoadsoih"
            ],
            "createdAt": "2021-01-01 00:00:00",
            "updatedAt": "2021-01-01 00:00:00"
        }
    },
    {
        method: 'DELETE',
        url: 'http://localhost:8080/memo', data: { wikiId: '21' },
        example: {
            "id": 21,
            "title": "askdaosjdoa",
            "content": "This is an introductory course to computer science",
            "ownerTraqId": "kavos",
            "tags": [
                "hoadsoih"
            ],
            "createdAt": "2021-01-01 00:00:00",
            "updatedAt": "2021-01-01 00:00:00"
        }
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/tag',
        example: [
            "ComputerScience"
        ]
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/wiki/user',
        example: [
            {
                "id": 123,
                "type": "sodan",
                "title": "Introduction to Computer Science",
                "Abstract": "This is an introductory course to computer science",
                "createdAt": "2021-01-01 00:00:00",
                "updatedAt": "2021-01-01 00:00:00",
                "ownerTraqId": "kavos",
                "tags": [
                    "ComputerScience"
                ]
            }
        ]
    },
    {
        method: 'POST',
        url: 'http://localhost:8080/wiki/user/favorite', data: { wikiId: '3' },
        example: {
            "id": 3,
            "type": "sodan",
            "title": "Introduction to Computer Science",
            "Abstract": "This is an introductory course to computer science",
            "createdAt": "2021-01-01 00:00:00",
            "updatedAt": "2021-01-01 00:00:00",
            "ownerTraqId": "kavos",
            "tags": [
                "ComputerScience"
            ]
        }
    },
    {
        method: 'GET',
        url: 'http://localhost:8080/wiki/user/favorite',
        example: [
            {
                "id": 3,
                "type": "sodan",
                "title": "Introduction to Computer Science",
                "Abstract": "This is an introductory course to computer science",
                "createdAt": "2021-01-01 00:00:00",
                "updatedAt": "2021-01-01 00:00:00",
                "ownerTraqId": "kavos",
                "tags": [
                    "ComputerScience"
                ]
            }
        ]
    },
    {
        method: 'DELETE',
        url: 'http://localhost:8080/wiki/user/favorite',
        data: { wikiId: '3' },
        example: {
            "id": 3,
            "type": "sodan",
            "title": "Introduction to Computer Science",
            "Abstract": "This is an introductory course to computer science",
            "createdAt": "2021-01-01 00:00:00",
            "updatedAt": "2021-01-01 00:00:00",
            "ownerTraqId": "kavos",
            "tags": [
                "ComputerScience"
            ]
        }
    }
];

async function executeRequests() {
    results = "";
    for (const req of requests) {
        try {
            const options = {
                method: req.method,
                headers: { 'Content-Type': 'application/json' },
                body: req.data ? JSON.stringify(req.data) : undefined
            };
            const response = await fetch(req.url, options);
            const data = await response.json();
            results += `===============${req.method} ${req.url}================\nresponse: ${JSON.stringify(data, null, "\t")}\nexpected: ${JSON.stringify(req.example, null, "\t")}\n\n`;
            //console.log(`Response from ${ req.method } ${ req.url }:`, data);
        } catch (error) {
            results += `===============${req.method} ${req.url}================\nERROR: ${error.message}\n\n`;
            //console.error(`Error in ${ req.method } ${ req.url }:`, error.message);
        }
    }

    fs.writeFileSync('results.txt', results);
}

executeRequests();