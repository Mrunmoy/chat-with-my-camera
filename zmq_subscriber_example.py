from publisher.zeromq_sub import ZeroMQSubscriber

def main():
    sub = ZeroMQSubscriber(port=5555)
    sub.subscribe()

if __name__ == "__main__":
    main()
