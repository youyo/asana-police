import asana
import datetime
import slack_sdk
import boto3


class AsanaPolice:
    def __init__(self):
        ssm_client = boto3.client("ssm")
        response = ssm_client.get_parameters(
            Names=[
                "/asana-police/ASANA_TOKEN",
                "/asana-police/ASANA_WORKSPACE_ID",
                "/asana-police/ASANA_PROJECT_ID",
                "/asana-police/SLACK_BOT_TOKEN",
                "/asana-police/SLACK_CHANNEL",
            ],
            WithDecryption=True,
        )
        params = {}
        for param in response["Parameters"]:
            params[param["Name"]] = param["Value"]

        self.client = asana.Client.access_token(params["/asana-police/ASANA_TOKEN"])
        self.project_id = params["/asana-police/ASANA_PROJECT_ID"]
        self.workspace_id = params["/asana-police/ASANA_WORKSPACE_ID"]
        self.slack_client = slack_sdk.WebClient(params["/asana-police/SLACK_BOT_TOKEN"])
        self.slack_channel = params["/asana-police/SLACK_CHANNEL"]

    def get_users(self):
        params = {
            "workspace": self.workspace_id,
        }
        users = self.client.users.get_users(params)
        return list(users)

    def get_expired_tasks(self, users: list):
        incomplete_tasks = []
        JST = datetime.timezone(datetime.timedelta(hours=+9), "JST")
        a_week_ago = datetime.datetime.now(JST) - datetime.timedelta(7)

        params = {
            "workspace": f"{self.workspace_id}",
            "completed_since": a_week_ago.isoformat(),
        }
        options = {
            "opt_fields": ["name", "due_on", "completed"],
        }

        for user in users:
            result = {
                "name": user["name"],
                "tasks": [],
            }

            params["assignee"] = user["gid"]

            tasks = self.client.tasks.get_tasks(params, **options)
            for task in tasks:
                if task["due_on"] is None:
                    continue
                if task["completed"] is True:
                    continue
                if (
                    datetime.datetime.strptime(task["due_on"], "%Y-%m-%d").replace(
                        tzinfo=JST
                    )
                    < a_week_ago
                ):
                    result["tasks"].append(
                        {
                            "name": task["name"],
                            "due_on": task["due_on"],
                            "url": f'https://app.asana.com/0/{self.project_id}/{task["gid"]}',
                        }
                    )

            incomplete_tasks.append(result)

        return incomplete_tasks

    def post_chat(self, messages):
        for msg in messages:
            self.slack_client.chat_postMessage(
                channel=self.slack_channel,
                text=msg,
            )


def handler(event, context):
    asana_police = AsanaPolice()
    users = asana_police.get_users()
    incomplete_tasks = asana_police.get_expired_tasks(users)

    messages = []

    for task in incomplete_tasks:
        if len(task["tasks"]) == 0:
            continue

        msg = f'{task["name"]}\n一週間以上期日の過ぎているタスクがあります。\n```\n'
        for t in task["tasks"]:
            msg = msg + f'{t["due_on"]} {t["name"]} {t["url"]}\n'
        msg = msg + "```"
        messages.append(msg)

    if len(messages) == 0:
        msg = "一週間以上期日の過ぎているタスクはありません🎉 \n"
        messages.append(msg)

    asana_police.post_chat(messages)
    return messages


if __name__ == "__main__":
    result = handler(None, None)
    import json

    print(json.dumps(result))
