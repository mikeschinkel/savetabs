# Roadmap
<var name="product" value="SaveTabs"/>

I have numerous plans for SaveTabs, documented below. However, the order listed is not indicative of the order in which I plan to add these features. That order will be determined by user interest.

Further, there is no specified timeline. Timeline will very much be influenced by my ability to find a way to monetize this effort, enough to at least cover what I would otherwise make as a freelance software developer doing contract work.

## Features 
Below are listed the features are would like to incorporate into SaveTabs:

### Import All The Bookmarks
A key part of my vision for SaveTabs is to provide a **single place that you can control** to import and maintain all your bookmarks, be they in Chrome, Edge, Safari, Opera, or any other browser, as well as from any bookmarking service online such as [Pocket](https://getpocket.com], [Raindrop](https://raindrop.io), [Pinboard](https://pinboard.in), [Diigo](https://www.diigo.com/) and more.   
#### Sync All the Bookmarks
While importing bookmarks is useful, ideally I will be able to implement one-way bookmark sync as well, so that anything new you add gets synced to your SaveTabs database.  

### Add Includes & Excludes
Allow adding patterns of URLs to include and/or URL patterns to exclude. For specifics, see [SaveTabs Include & Exclude Patterns](SaveTabs-Includes-Excludes.md) 
   
### Alternate Daemons
As an option to using my locally-installed dameon on your own computer I plan to implement other versions of the daemon that can run in a serverless environment using AWS Lambda and other similar cloud services. This would allow 
#### Serverless daemon
I envision offering a daemon that AWS Lambda and other similar cloud services you could host on your own AWS or other cloud account.
#### Server-based daemon
I envision offering a daemon that can be hosted on any Linux server such as ones provided by Digital Ocean, or practically any other provider  
#### NAS-based daemon
In addition to external support I envision offering daemon packaged to run on popular NAS products such as Synology, QNAP, TrueNAS, and more.  

### Daemon Security
At this time the daemon provides a RESTful API with no security in place since SaveTabs only current option is a local install, and I are trying to get an [MVP](https://en.wikipedia.org/wiki/Minimum_viable_product) developed. However, in the future I envision securing the daemon in a potential variety of ways no matter where is it hosted.
#### API Key Security
Initially I envision allowing your daemon to secured by an API key, but I do need to research security best practices before making commitments on specifically how that would be implemented.
#### OAuth Security
I envision allowing your daemon to secured by any authentication provider such as GitHub, Google, LinkedIn, Facebook, and others.
#### 2-factor Auth
Further, I envision adding support for 2-factor authentication at some point.

### Alternate Database Support
As an option to only supporting storage of data in a Sqlite database, I envision adding support for other databases. 
#### SQL databases
**Postgres** support is a definite, and MySQL support may also be in the cards. Support for other databases will be determined by user demand.
#### NoSQL databases
MongoDB support may be an option, as well as potentially other NoSQL databases.

### Multi-user/Team/Family Support
If SaveTabs becomes popular enough that multiple people such as teams and/or families want to be able to share an account and share a daemon I will consider adding such functionality.  This might be one way in which I could monetize, by providing a base multi-user version as open-source, with additional paid functionality for teams and/or families. 

### LLM-powered AI Chatbot
A key goal for SaveTabs is to provide a user-specific corpus of web resources that can be trained for a given user on all the content they have consumed on that web since they have started bookmarking and tagging content on the web.

### Search Engine
Similar to my AI Chatbot goal, I want SaveTabs to be able to provide a super-memory for a user of all the things they have seen on the web. How this will be implemented and what it will look like remains to be seen, however. 

# Feedback
If you are interested in seeing support for any specific type of feature, or may even want to sponsor SaveTabs and/or specific supported features, please reach out to me [via email](mailto:info@gearbox.works). 
