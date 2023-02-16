## Privacy: where is your data?

Hello! I am [Xavier Vello](https://github.com/xvello), a French
[software engineer](https://linkedin.com/in/xaviervello) and maintainer of this project.

The project does not have a legal team to write a proper privacy policy, but please read the following
information and email me at `hello@<this-domain>.it` if you have any question.

### Identity management

To create an account on the website, a valid e-mail address and a password are required.
This address will only be used for account recovery and to contact you if unusual activity is detected.
I will **never** send unsolicited newsletters or share your e-mail with other third-parties than ORY Corp.

Because I don't trust myself with your personal information, I am delegating this to professionals: the
[Ory Cloud](https://ory.sh/docs) service is used to store and secure your credentials. Here is
[their privacy policy](https://www.ory.sh/privacy/). Their knowledge of your activity on this website
is limited to your sessions' technical data (at what times you are active, and from what device / IP),
they don't know or care about what pages you see, or how you interact with them.

### Main database and servers

The main database is hosted on [Hetzner Cloud](https://www.hetzner.com/cloud), in their Nuremberg (Germany) datacenter.
Servers are secured to the best of my knowledge and abilities, and the source of this website is
[available on GitHub](https://github.com/letsblockit/letsblockit), under the Apache License version 2.0.

The applicative servers are running on the [Fly.io platform](https://fly.io/), both in Europe and the USA.
Their CDN will redirect your requests to the closest server.

Thanks to delegating the authentication to Ory, the servers don't know or store your e-mail or password, just
[a random unique identifier](https://en.wikipedia.org/wiki/Universally_unique_identifier). This means that
a leak of the database would not compromise your credentials.

While aggregated usage metrics will be computed (how many users use a given filters, with how many parameters),
the parameter values will **never** be accessed outside of maintenance operation (database recovery, data model
migrations). Although it would be pretty valuable to extract new filters, your privacy is more important. Please
[suggest new filters to help the project](/help/contributing) instead of keeping them as custom rules!

### Warning: filter lists are downloadable without authentication

Because ad-blockers are designed to use public blocking lists, they don't support authenticating when downloading a
list. Here is how your list is secured:

- Your list is assigned [a random unique identifier](https://en.wikipedia.org/wiki/Universally_unique_identifier),
separate from your user account. This ID acts as a private download token: your list will be accessible on
`https://get.letsblock.it/list/$token`
- If you think the URL has been leaked, you can generate a new token from your [account settings](/user/account) page.
This will block any download via the old URL.
- The random ID reduces the risk of enumeration attacks. Other protections may be present, but I'll keep
some obscurity here to keep them effective :)
- In the future, I plan on allowing users to maintain several lists, to share one while keeping some filters private.
This will also allow users to split their filters into several private lists to avoid one single list holding enough
information to identify them. Contributions are welcome to move this forward if you are interested!

### Zero third-party tracking

While some javascript is present for progressive enhancement (shout-out to the [htmx](https://htmx.org/) project),
third-party tracking will **never** be present on this website. Access logs are indexed and analysed to build
aggregated usage metrics and to combat abuse.
