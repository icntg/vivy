import asyncio
import aiomysql

# 返回一个asyncio事件循环。
loop = asyncio.get_event_loop()

async def test_example():
    pool = await aiomysql.create_pool(
        host='127.0.0.1',
        port=3306, user='root',
        password='root',
        loop=loop,
        cursorclass=aiomysql.DictCursor
    )

    try:
        async with pool.acquire() as conn:
            async with conn.cursor(aiomysql.SSCursor) as cur:
                # await cur.execute('''CREATE DATABASE `vivy` /*!40100 COLLATE 'utf8mb4_bin' */;''')
                await cur.execute('''CREATE DATABASE IF NOT EXISTS `vivy` /*!40100 COLLATE 'utf8mb4_bin' */;''')
                await cur.execute('''CREATE USER IF NOT EXISTS 'vivy'@'%' IDENTIFIED BY 'vivy';''')
                await cur.execute('''GRANT ALL PRIVILEGES ON vivy.* TO 'vivy'@'%';''')
    except Exception as e:
        print(e)
    pool.close()


loop.run_until_complete(test_example())
